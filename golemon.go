package main

import "fmt"
// import "strconv"
import "strings"
import "time"
import "os"
import "os/signal"
import "syscall"
import "io/ioutil"
import "os/exec"

const separator_left string = ""
const separator_right string = ""
const color_first string = "#FF282A2E"
const color_second string = "#FF454A4F"

func initBlocks() map[string]string {
    var b map[string]string
    b = make( map[string]string )
    b["music"] = ""
    b["workspaces"] = ""
    b["torrent"] = ""
    b["volume"] = ""
    // b["cpu"] = ""
    b["battery"] = ""
    b["brightness"] = ""
    b["redshift"] = ""
    b["wifi"] = ""
    b["layout"] = ""
    b["date"] = ""
    return b
}
func refreshBlock(block string) string{
    block = strings.TrimSuffix(block, "\n")
    var result string
    switch block{
    case "music":
        musicScript := exec.Command("/home/infiniter/Code/GoLemon/music")
        musicScriptOut, _ := musicScript.Output()
        result = "%{B" + color_first + "}" + string(musicScriptOut) + " %{B-}%{F" + color_first + "}" + separator_right + "%{F-}"
    case "workspaces":
        workspacesScript := exec.Command("/home/infiniter/Code/GoLemon/workspaces")
        workspacesScriptOut, _ := workspacesScript.Output()
        result = string(workspacesScriptOut)
    case "torrent":
        torrentScript := exec.Command("/home/infiniter/Code/GoLemon/torrent")
        torrentScriptOut, _ := torrentScript.Output()
        result = "%{F" + color_first + "}" + separator_left + "%{F-}%{B" + color_first + "} " + string(torrentScriptOut) + " %{F" + color_second + "}" + separator_left + "%{F-}%{B-}"
    case "volume":
        volumeScript := exec.Command("/home/infiniter/Code/GoLemon/volume")
        volumeScriptOut, _ := volumeScript.Output()
        result = "%{B" + color_second + "} " + string(volumeScriptOut) + " %{F" + color_first + "}" + separator_left + "%{F-}%{B-}"
    case "battery":
        batteryScript := exec.Command("/home/infiniter/Code/GoLemon/battery")
        batteryScriptOut, _ := batteryScript.Output()
        result = "%{B" + color_first + "} " + string(batteryScriptOut) + " %{F" + color_second + "}" + separator_left + "%{F-}%{B-}"
    case "brightness":
        brightnessScript := exec.Command("/home/infiniter/Code/GoLemon/brightness")
        brightnessScriptOut, _ := brightnessScript.Output()
        result = "%{B" + color_second + "} " + string(brightnessScriptOut) + " %{F" + color_first + "}" + separator_left + "%{F-}%{B-}"
    case "redshift":
        redshiftScript := exec.Command("/home/infiniter/Code/GoLemon/redshift")
        redshiftScriptOut, _ := redshiftScript.Output()
        result = "%{B" + color_first + "} " + string(redshiftScriptOut) + " %{F" + color_second + "}" + separator_left + "%{F-}%{B-}"
    case "wifi":
        wifiScript := exec.Command("/home/infiniter/Code/GoLemon/wifi")
        wifiScriptOut, _ := wifiScript.Output()
        result = "%{B" + color_second + "} " + string(wifiScriptOut) + " %{F" + color_first + "}" + separator_left + "%{F-}%{B-}"
    case "layout":
        layoutScript := exec.Command("/home/infiniter/Code/GoLemon/layout")
        layoutScriptOut, _ := layoutScript.Output()
        result = "%{B" + color_first + "} " + string(layoutScriptOut) + " %{F" + color_second + "}" + separator_left + "%{F-}%{B-}"
    case "date":
        dateScript := exec.Command("/home/infiniter/Code/GoLemon/date")
        dateScriptOut, _ := dateScript.Output()
        result = "%{B" + color_second + "} " + string(dateScriptOut) + " %{B-}"
    }
    return result
}
func prepareForLemon( blocks map[string]string ) string{
    left := []string{"music"}
    center := []string{"workspaces"}
    right := []string{"torrent", "volume", "battery", "brightness", "redshift", "wifi", "layout", "date"}
    vals_left := make([]string, 0, len(left))
    vals_center := make([]string, 0, len(center))
    vals_right := make([]string, 0, len(right))
    for _,key_left := range left{
        vals_left = append(vals_left, blocks[key_left])
    }
    for _,key_center := range center{
        vals_center = append(vals_center, blocks[key_center])
    }
    for _,key_right := range right{
        vals_right = append(vals_right, blocks[key_right])
    }
    res := "%{l}" + strings.Join(vals_left, "") + "%{c}" + strings.Join(vals_center, "") + "%{r}" + strings.Join(vals_right, "")
    return res
}
func fetchEmpty( data map[string] string ) map[string] string{
    for key,val := range data{
        if val == "" {
            data[key] = refreshBlock( key )
        }
    }
    return data
}
func main() {
    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)

    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)

    fmt.Println("Starting GoLemon")
    go func() {
        lemon := exec.Command("lemonbar", "-p", "-f", "Inconsolata for Powerline-15", "-f", "FontAwesome-16", "-B", "#0b0b0b", "-F", "#CCCCCC", "-g", "1920x25+0+0")
        lemonIn, _ := lemon.StdinPipe()
        // lemonOut, _ := lemon.StdoutPipe()
        lemon.Start()

        stalonetray := exec.Command("stalonetray" ,"--geometry", "3x1+700+0", "--grow-gravity", "W", "--icon-gravity", "E", "-bg", "#0b0b0b", "--max-geometry", "10x1")
        // stalonetrayIn, _ := stalonetray.StdinPipe()
        // stalonetrayOut, _ := stalonetray.StdoutPipe()
        stalonetray.Start()

        go func() {
            time.Sleep(time.Second * 3)
            bar := exec.Command("fix_layers_golemon")
            // bar_id, _ := bar.StdoutPipe()
            bar.Start()
        }()
        var blocks map[string]string
        blocks = initBlocks()
        blocks = fetchEmpty(blocks)
        for {
            lemonIn.Write([]byte(prepareForLemon(blocks)))
            sig := <-sigs
            switch sig {
            case syscall.SIGINT:
                fmt.Println("interrupted")
                done <- true
            case syscall.SIGTERM:
                fmt.Println("terminated")
                done <- true
            case syscall.SIGUSR1:
                dat, _ := ioutil.ReadFile("/tmp/golemon_refresh")
                blocks[strings.TrimSuffix(string(dat), "\n")] = refreshBlock(string(dat))
            }
        }

        done <- true
        lemonIn.Close()
    }()

    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}
