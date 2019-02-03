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

const separator string = " %{F#c22330}•%{F-} "

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
        result = string(musicScriptOut)
    case "workspaces":
        workspacesScript := exec.Command("/home/infiniter/Code/GoLemon/workspaces")
        workspacesScriptOut, _ := workspacesScript.Output()
        result = string(workspacesScriptOut)
    case "torrent":
        torrentScript := exec.Command("/home/infiniter/Code/GoLemon/torrent")
        torrentScriptOut, _ := torrentScript.Output()
        result = string(torrentScriptOut)
    case "volume":
        volumeScript := exec.Command("/home/infiniter/Code/GoLemon/volume")
        volumeScriptOut, _ := volumeScript.Output()
        result = string(volumeScriptOut)
    // case "cpu":
    //     cpuScript := exec.Command("/home/infiniter/Code/GoLemon/cpu")
    //     cpuScriptOut, _ := cpuScript.Output()
    //     result = string(cpuScriptOut)
    case "battery":
        batteryScript := exec.Command("/home/infiniter/Code/GoLemon/battery")
        batteryScriptOut, _ := batteryScript.Output()
        result = string(batteryScriptOut)
    case "brightness":
        brightnessScript := exec.Command("/home/infiniter/Code/GoLemon/brightness")
        brightnessScriptOut, _ := brightnessScript.Output()
        result = string(brightnessScriptOut)
    case "redshift":
        redshiftScript := exec.Command("/home/infiniter/Code/GoLemon/redshift")
        redshiftScriptOut, _ := redshiftScript.Output()
        result = string(redshiftScriptOut)
    case "wifi":
        wifiScript := exec.Command("/home/infiniter/Code/GoLemon/wifi")
        wifiScriptOut, _ := wifiScript.Output()
        result = string(wifiScriptOut)
    case "layout":
        layoutScript := exec.Command("/home/infiniter/Code/GoLemon/layout")
        layoutScriptOut, _ := layoutScript.Output()
        result = string(layoutScriptOut)
    case "date":
        dateScript := exec.Command("/home/infiniter/Code/GoLemon/date")
        dateScriptOut, _ := dateScript.Output()
        result = string(dateScriptOut)
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
    res := "%{l}" + strings.Join(vals_left, separator) + "%{c}" + strings.Join(vals_center, separator) + "%{r}" + strings.Join(vals_right, separator)
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
        lemon := exec.Command("lemonbar", "-p", "-f", "League Mono-14", "-f", "FontAwesome-16", "-B", "#000000", "-F", "#CCCCCC", "-g", "1920x25+0+0", "| bash")
        lemonIn, _ := lemon.StdinPipe()
        // lemonOut, _ := lemon.StdoutPipe()

        lemon.Start()

        go func() {
            time.Sleep(time.Second * 3)
            bar := exec.Command("fix_layers_golemon")
            bar_id, _ := bar.StdoutPipe()
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
