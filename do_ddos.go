package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/yahoo/vssh"
	"log"
	"os"
	"strings"
	"time"
)

// I will be hardcode these machine names because it requires a lot more machine power to generate random names
var bulkNames = [10]string{"m1", "m2", "m3", "m4", "m5", "m6", "m7", "m8", "m9", "m10"}

func main() {
	// get api key from enviroment variable
	token := os.Getenv("DIGITALOCEAN_TOKEN")

	// authentication
	client := godo.NewFromToken(token)
	// avaliable regions
	var regions = [9]string{"nyc1", "nyc3", "sfo3", "ams3", "sgp1", "lon1", "fra1", "tor1", "blr1"} // ams3 will be default value for region if not defined otherwise
	var size = [5]string{"s-1vcpu-1gb", "s-1vcpu-2gb", "s-2vcpu-2gb", "s-2vcpu-4gb", "s-4vcpu-8gb"} // s-1vcpu-1gb will be default value for size if not defined otherwise
	image := "ubuntu-20-04-x64"                                                                     // default value for image is 20.04 ubuntu, others will be added later

	// flags are defined below
	// Create single droplet
	csdPtr := flag.Bool("csd", false, "Create single droplet with only required parameters.")
	// Create single droplet with tag
	csdwtPtr := flag.Bool("csdwt", false, "Create single droplet with tags.")
	// Create multiple droplets with tag
	cmdwtPtr := flag.Bool("cmdwt", false, "Create multiple droplets with tag, droplet names are created incremantally.")
	// delete single droplet with ID
	ddwiPtr := flag.Int("ddwi", 0, "A non-zero integer value which belongs to the droplet that wanted to be deleted.")
	// delete droplets with tag name
	ddwtPtr := flag.String("ddwt", "", "String value belongs which belongs to the droplet(s) that wanted to be deleted.")
	// define a droplet name
	dropletName := flag.String("dname", "", "Droplet name, required for creating a single droplet.")
	// define a tag
	dtagPtr := flag.String("dtag", "", "Droplet tag, can be used for indicate a group of droplets.")
	// list all droplets, also works with dtag
	ladPtr := flag.Bool("lad", false, "List all droplets. If this flag gets combined with dtag flag then only droplets with provided tag name will be listed.")
	// define amount of droplets to be created, required to used with cmdwt
	amountPtr := flag.Int("amount", 0, "Amount of droplets that will be created, must be used with th cmdwt flag")
	// define a ssh fingerprint to be assigned to created droplet(s), check notes.txt to learn how to get one
	sshfpPtr := flag.String("sshfp", "", "SSH fingerprint, needed for connecting the droplets, and make them do things")
	// define a command to be run on all currently active machines
	rcPtr := flag.String(`rc`, "", "Command to be run on all active , should be used with dtag flag to specify the tag")
	flag.Parse()

	if *csdPtr {
		// cannot retrieve IP addr of the droplet immediately after creation request, it takes time to create a droplet on background
		droplet, _, err := createSingleDroplet(*dropletName, regions[3], size[1], image, *sshfpPtr, "", client)
		if err != nil {
			fmt.Printf("csd err : %s\n", err)
		}
		//ipaddr , err := droplet.PublicIPv4()
		ipaddr := droplet.Networks.V4[0].IPAddress

		if err != nil {
			fmt.Printf("err when getting ipv4 address : %s", err)
		}

		fmt.Printf("Droplet created succesfully \n ID: %d\n IP Addr: %s", droplet.ID, ipaddr)
	}
	///////////////////
	if *csdwtPtr {
		droplet, _, err := createSingleDroplet(*dropletName, regions[3], size[1], image, *sshfpPtr, *dtagPtr, client)
		if err != nil {
			fmt.Printf("csd err : %s\n", err)
		}
		ipaddr := droplet.Networks.V4[0].IPAddress

		if err != nil {
			fmt.Printf("err when getting ipv4 address : %s", err)
		}

		fmt.Printf("Droplet created succesfully \n ID: %d\n Tag(s): %s\n IP Addr: %s", droplet.ID, droplet.Tags, ipaddr)

	}
	///////////////////
	if *cmdwtPtr {
		droplets, _, err := createMultipleDroplets(regions[3], size[1], image, *sshfpPtr, *dtagPtr, *amountPtr, client)

		if err != nil {
			fmt.Printf("cmdwt err : %s\n", err)
		}
		getDropletInfo(droplets)
	}
	///////////////////
	if *rcPtr != "" {


		IPAddrs, err := getIPAddressses(client, *dtagPtr)
		if err != nil {
			fmt.Printf("Error while retrieving IP addresses : %v/n", err)
			return
		}

		err = runCommandOnGivenIPAddresses(*rcPtr, IPAddrs)

		if err != nil {
			fmt.Printf("Error while running commands on given machines : %v/n", err)
			return
		}
		// SSH TEST ZONE
	}
	///////////////////
	if *ddwiPtr != 0 {
		err := deleteDroplet(client, *ddwiPtr)
		if err != nil {
			fmt.Printf("err when deleting a droplet : %s", err)
		}
		fmt.Printf("Droplet deleted succesfully\n")
	}
	///////////////////
	if *ddwtPtr != "" {
		err := deleteDroplet(client, *ddwtPtr)
		if err != nil {
			fmt.Printf("err when deleting a droplet : %s", err)
		}
		fmt.Printf("Droplet deleted succesfully\n")
	}
	///////////////////
	if *ladPtr {
		droplets, err := listDroplets(client, *dtagPtr)
		if err != nil {
			fmt.Printf("err when listing the droplets: %s", err)
		}
		getDropletInfo(droplets)
	}
}

// createSingleDroplet creates a single droplet , if user provides a tag value, the droplet will be created accordingly
func createSingleDroplet(name, region, size, image, sshfp, tag string, client *godo.Client) (*godo.Droplet, *godo.Response, error) {
	ctx := context.TODO()

	//  Hacky solution for testing purposses I will fix that later
	var tags = []string{tag}
	if tag == "" {
		tags = nil
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: image,
		},
		Tags: tags,
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: sshfp},
		},
	}

	newdroplet, resp, err := client.Droplets.Create(ctx, createRequest)

	if err != nil {
		return nil, nil, err
	}

	readyDroplet, isCompleted, err := waitForDropletToBeCompleted(ctx, client, newdroplet.ID)

	if err != nil {
		return nil, nil, err
	}

	if !isCompleted {
		return nil, nil, err
	} else {
		return readyDroplet, resp, nil
	}

}

// createMultipleDroplets creates multiple droplets, assigns user provided tag to all of them to make it easier to control them
func createMultipleDroplets(region, size, image, sshfp, tag string, amount int, client *godo.Client) ([]godo.Droplet, *godo.Response, error) {
	ctx := context.TODO()
	var readyDroplets []godo.Droplet
	//  Hacky solution for testing purposses I will fix that later
	var tags = []string{tag}
	if tag == "" {
		tags = nil
	}

	createRequest := &godo.DropletMultiCreateRequest{

		Names:  bulkNames[:amount],
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: image,
		},
		Tags: tags,
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: sshfp},
		},
	}

	newDroplets, resp, err := client.Droplets.CreateMultiple(ctx, createRequest)

	if err != nil {
		return nil, nil, err
	}
	for _, droplet := range newDroplets {
		readyDroplet, isCompleted, err := waitForDropletToBeCompleted(ctx, client, droplet.ID)

		if err != nil {
			return nil, nil, err
		}

		if !isCompleted {
			return nil, nil, err
		} else {
			readyDroplets = append(readyDroplets, *readyDroplet)
		}
	}

	return readyDroplets, resp, nil
}

// deleteDroplet deletes the droplet whether it specified with its id or tag
// if tag is provided then function deletes all droplets that have the same tag
func deleteDroplet(client *godo.Client, indicator interface{}) (err error) {
	ctx := context.TODO()

	// Type Switch for less code
	switch ind := indicator.(type) {
	case int:
		_, err := client.Droplets.Delete(ctx, ind)
		return err

	case string:
		_, err := client.Droplets.DeleteByTag(ctx, ind)
		return err
	default:
		return errors.New("Geçersiz Droplet tanımlayıcısı: ID(integer) veya TAG(String) girilmeli. ")
	}

}

// waitForDropletToCompleted checks the given droplet for every 5 seconds and return true if its completed
// and returns the droplet with ready state
func waitForDropletToBeCompleted(ctx context.Context, client *godo.Client, id int) (*godo.Droplet, bool, error) {
	spinner := [5]string{"          ", "==        ", "=====     ", "========  ", "=========="}
	counter := 0
	tick := time.Tick(500 * time.Millisecond)
	timeout := time.After(60 * time.Second)
	for {
		select {
		case <-timeout:
			return nil, false, errors.New("timeout waiting for droplet")
		case <-tick:
			counter = (counter + 1) % len(spinner)
			fmt.Printf("\r[%s] Waiting for state: %s", spinner[counter], "Created")
			droplet, _, _ := client.Droplets.Get(ctx, id)
			if droplet.Status == "active" {
				return droplet, true, nil
			}
		}
	}
}

// listDroplets list all droplets if there is no tag provided otherwise it returns droplets which has given tag
// I will divide this function in two parts to prevent more comparison
func listDroplets(client *godo.Client, tag string) ([]godo.Droplet, error) {
	ctx := context.TODO()

	if tag == "" {
		returnedDroplets, _, err := client.Droplets.List(ctx, nil)

		if err != nil {
			return nil, err
		} else {
			return returnedDroplets, nil
		}

	} else {
		returnedDroplets, _, err := client.Droplets.ListByTag(ctx, tag, nil)
		if err != nil {
			return nil, err
		} else {
			return returnedDroplets, nil
		}
	}
}

// getIPAddressses uses listDroplets to retrieve IP addresses of already created droplets
func getIPAddressses(client *godo.Client, tag string) ([]string, error) {
	droplets, err := listDroplets(client, tag)

	if err != nil {
		return nil, err
	}
	ipaddrs := make([]string, len(droplets))

	// get list of all machine's IP addresses
	for i, droplet := range droplets {
		ipaddrs[i] = droplet.Networks.V4[0].IPAddress
	}

	return ipaddrs, nil
}

func getDropletInfo(droplets []godo.Droplet) {
	separator := strings.Repeat("*", 26)
	fmt.Printf("\n%s", separator)
	for _, droplet := range droplets {
		fmt.Printf("\nDroplet ID: %d\nIPAddress: %s\nTag(s): %s\n%s\n", droplet.ID, droplet.Networks.V4[0].IPAddress, droplet.Tags, separator)
	}
}

// runCommandOnGivenIPAddresses runs given command on the machines concurrently with the help of VSSH package
func runCommandOnGivenIPAddresses(cmd string, Addresses []string) error {

	vs := vssh.New().Start()
	// droplets doesn't have any ssh password setted
	config, err := vssh.GetConfigPEM("root", "C:/Users/yusuf/.ssh/id_rsa")

	if err != nil {
		return err
	}

	for _, addr := range Addresses {
		err := vs.AddClient(addr+":22", config, vssh.SetMaxSessions(1))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Client succesfully added")
		}

	}
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	timeout, _ := time.ParseDuration("100s")

	respchan := vs.Run(ctx, cmd, timeout)

	for resp := range respchan {
		if err := resp.Err(); err != nil {
			log.Println(err)
			continue
		}

		outTxt, errTxt, _ := resp.GetText(vs)
		fmt.Println(outTxt, errTxt, resp.ExitStatus())

	}
	return nil
}

// TODO: Be able to create and boot machines through digital ocean API-------------------------------------------------------------> DONE
// TODO: Creating can be done with OR without tags---------------------------------------------------------------------------------> DONE
// TODO: Deleting can be done with ID or TAG---------------------------------------------------------------------------------------> DONE
// TODO: Be able to delete machine-------------------------------------------------------------------------------------------------> DONE
// TODO: Show machine info (IP, name, id etc.)-------------------------------------------------------------------------------------> DONE
// TODO: Be able to list all droplets or droplets that has same tag----------------------------------------------------------------> DONE
// TODO: Be able to create multiple droplets at once (bind them with a tag)--------------------------------------------------------> DONE
// TODO: Be able to define a SSH key to created machine----------------------------------------------------------------------------> DONE
// TODO: Be able to restart machine------------------------------------------------------------------------------------------------> DONE
// TODO: Be able to operate created machines via ssh (make them send requests to the target etc.)----------------------------------> DONE
// TODO: Add this project to github with a proper documentation and add "regrews" as a contributer IMPORTANT
// TODO: Add a killswitch (commad is "killall -u root")  IMPORTANT
// TODO: Be able to restart machines with tag name		 IMPORTANT
// TODO: Be able define ssh key inside of the code		 IMPORTANT
// TODO: Be able to run scripts on created machines		 IMPORTANT
// TODO: Be able to see only one machine's output with flag  'logging mechanism' XML
// TODO: Be able to see all machine's outputs with another flag 'logging mechanism' XML
// TODO: Try to make a WEB interface with react
// TODO: Carry this completed program to "COBRA" framework
// TODO: Log anomalies happened during execution