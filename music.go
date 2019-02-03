package main

import (
    "fmt"
    "strings"
    "os/exec"
    "regexp"
    // "strconv"
)

func main(){
    var result string
    mpcCommand := exec.Command("mpc")
    mpcOut, _ := mpcCommand.Output()
    mpcArray := strings.Split(string(mpcOut), "\n")
    if len(mpcArray) < 3 {
        result = "%{F#555555}  playlist empty%{F-}"
    } else {
        name := mpcArray[0]
        data := mpcArray[1]
        flags := mpcArray[2]

        if len(name) >= 45 {
            name = name[:41] + "..." + name[len(name)-3:]
        }

        // FLAGS
        // r_volume := regexp.MustCompile("volume:([0-9]+)")
        // capture_volume := r_volume.FindStringSubmatch(flags)
        // volume, _ := strconv.Atoi(capture_volume[1])
        // var volume_res string
        // if volume < 100  {
        //     volume_res = "- " + string(volume) + "%"
        // } else {
        //     volume_res = ""
        // }

        r_repeat := regexp.MustCompile("repeat: ([a-z]+)")
        capture_repeat := r_repeat.FindStringSubmatch(flags)
        repeat := capture_repeat[1]
        var repeat_res string
        if repeat == "on" {
            repeat_res = ""
        } else {
            repeat_res = "%{F#555555}%{F-}"
        }

        r_random := regexp.MustCompile("random: ([a-z]+)")
        capture_random := r_random.FindStringSubmatch(flags)
        random := capture_random[1]
        var random_res string
        if random == "on" {
            random_res = ""
        } else {
            random_res = "%{F#555555}%{F-}"
        }

        r_single := regexp.MustCompile("single: ([a-z]+)")
        capture_single := r_single.FindStringSubmatch(flags)
        single := capture_single[1]
        var single_res string
        if single == "on" {
            single_res = ""
        } else {
            single_res = "%{F#555555}%{F-}"
        }

        r_consume := regexp.MustCompile("consume: ([a-z]+)")
        capture_consume := r_consume.FindStringSubmatch(flags)
        consume := capture_consume[1]
        var consume_res string
        if consume == "on" {
            consume_res = ""
        } else {
            consume_res = "%{F#555555}%{F-}"
        }

        flags_res := " " + random_res + " " + consume_res + " " + single_res + " " + repeat_res

        // PLAYING
        r_playing := regexp.MustCompile("\\[([a-z]+)\\]")
        capture_playing := r_playing.FindStringSubmatch(data)
        playing := capture_playing[1]

        if playing == "playing" {
            result = "  " + name + flags_res
        } else {
            result = " %{F#555555} " + name + "%{F-}" + flags_res
        }

    }
    fmt.Print(result)
}
