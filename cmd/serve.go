/*
Copyright © 2020 Nick Maliwacki <knic.knic@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a tcp port and forwards that to a unix domain socket",
	Long: `Serves a tcp port and forwards that to a unix domain socket

examples:
    Bind only to localhost
        serve --path .\socket --address localhost:8080
    Bind to all ips
        serve --path .\socket --address :8080
        `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		server, err := net.Listen("tcp", Address)
		if err != nil {
			panic(err)
		}
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			go serveConnection(conn)
		}
	},
}

func serveConnection(conn net.Conn) {
	defer conn.Close()

	domain, err := net.ResolveUnixAddr("unix", Path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	remote, err := net.DialUnix("unix", nil, domain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer remote.Close()
	proxy(remote, conn)
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
