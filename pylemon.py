import atexit
import signal
import subprocess
import os
import time

mutex = False

class Pylemon(object):
    def __init__(self):
        self.outputs = dict()
        self.outputs['left'] = dict()
        self.outputs['right'] = dict()
        self.outputs['center'] = dict()
        self.first = True
        self.separator = '  %{F#c22330}•%{F-}  '

        # sets order of modules
        self.states = {
                'torrent': False,
                'volume': False,
                'cpu': False,
                'battery': False,
                'brightness': False,
                'redshift': False,
                'wifi': False,
                'layout': False,
                'date': False,
                'music': False,
                'workspaces': False
        }

        self.functions = {
                'torrent': self.get_torrent,
                'wifi': self.get_wifi,
                'volume': self.get_volume,
                'brightness': self.get_brightness,
                'cpu': self.get_cpu,
                'battery': self.get_battery,
                'date': self.get_date,
                'layout': self.get_layout,
                'redshift': self.get_redshift,
                'music': self.get_music,
                'workspaces': self.get_workspaces
        }

        self.positions = {
                'torrent': 'right',
                'wifi': 'right',
                'volume': 'right',
                'brightness': 'right',
                'cpu': 'right',
                'battery': 'right',
                'layout': 'right',
                'date': 'right',
                'redshift': 'right',
                'music': 'left',
                'workspaces': 'center'
        }

        self.lemonbar = subprocess.Popen(['lemonbar', '-p', '-f', '"League Mono"-14', '-f', 'FontAwesome-16', '-B', '#000000', '-F', '#CCCCCC', '-g', '1920x25+0+0'], stdin=subprocess.PIPE, stdout=subprocess.PIPE)

        subscribe_cpu = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_cpu', '3'])
        self.subscribe_cpu_pid = subscribe_cpu.pid

        subscribe_date = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_date', '60'])
        self.subscribe_date_pid = subscribe_date.pid

        subscribe_battery = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_battery', '60'])
        self.subscribe_battery_pid = subscribe_battery.pid

        sub_workspace = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_workspaces'])
        self.sub_workspace_pid  = sub_workspace.pid

        sub_music = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_music'])
        self.sub_music_pid  = sub_music.pid

        # sub_volume = subprocess.Popen(['/home/infiniter/Code/Pylemon/subscribe_volume'])
        # self.sub_volume_pid  = sub_volume.pid

        stalonetray = subprocess.Popen(['stalonetray' ,'--geometry', '5x1+680+0', '--grow-gravity', 'W', '--icon-gravity', 'E', '-bg', '#000000', '--max-geometry', '10x1'])
        self.stalonetray_pid  = stalonetray.pid

        self.executer = subprocess.Popen(['bash'], stdin=self.lemonbar.stdout)

        self.run()

    def get_date(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/date'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_layout(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/layout'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_redshift(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/redshift'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_music(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/music'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_cpu(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/cpu'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_battery(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/battery'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_brightness(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/brightness'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_workspaces(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/workspaces'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_volume(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/volume'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_wifi(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/wifi'], stdout=subprocess.PIPE)
        return result.stdout.decode()
    def get_torrent(self):
        result = subprocess.run(['/home/infiniter/Code/Pylemon/torrent'], stdout=subprocess.PIPE)
        return result.stdout.decode()

    def refresh_user(self, *args, **kwargs):
        global mutex
        if mutex:
            return
        mutex = True

        try:
            with open('/tmp/refresh', 'r') as t:
                target = t.read().replace('\n','')
            with open('/tmp/refresh', 'w') as t:
                t.write('')
        except FileNotFoundError:
            with open('/tmp/refresh', 'w') as t:
                t.write('')
            return
        if target == 'date':
            self.states['date'] = False
        elif target == 'brightness':
            self.states['brightness'] = False
        elif target == 'wifi':
            self.states['wifi'] = False
        elif target == 'redshift':
            self.states['redshift'] = False
        elif target == 'music':
            self.states['music'] = False
        elif target == 'layout':
            self.states['layout'] = False
        elif target == 'workspaces':
            self.states['workspaces'] = False
        elif target == 'cpu':
            self.states['cpu'] = False
        elif target == 'date':
            self.states['date'] = False
        elif target == 'volume':
            self.states['volume'] = False
        elif target == 'torrent':
            self.states['torrent'] = False
        elif target == 'battery':
            self.states['battery'] = False
        self.refresh()
        mutex = False
    def refresh_timer(self, *args, **kwargs):
        for key in self.states:
            self.states[key] = False
        self.refresh()

    def refresh(self):
        for key in self.states:
            if self.states[key] is False:
                self.outputs[self.positions[key]][key] = self.functions[key]()
                self.states[key] = True
        left = '%{l}' + self.separator.join([x for x in self.outputs['left'].values() if x != ''])
        center = '%{c}' + self.separator.join([x for x in self.outputs['center'].values() if x != ''])
        right = '%{r}' + self.separator.join([x for x in self.outputs['right'].values() if x != ''])
        self.lemonbar.stdin.write(' {} '.format(left + center + right).encode())
        self.lemonbar.stdin.flush()
        if self.first:
            os.system('xdo above -t "$(xdo id -N Bspwm -n root | sort | head -n 1)" "$(xdo id -a bar)"')
            os.system('xdo above -t "$(xdo id -a bar)" "$(xdo id -a stalonetray)"')
            self.first = False

    def run(self):
        signal.signal(signal.SIGSEGV, self.refresh_timer)
        signal.signal(signal.SIGUSR1, self.refresh_user)
        signal.signal(signal.SIGPIPE, signal.SIG_DFL)
        # initial paint
        self.refresh()
        while True:
            signal.pause()

    def kill_child_processes():
        print('cleaning up child processes')
        # subprocess.Popen(['pkill', '-f', 'pylemon_wakeup'])
        subprocess.Popen(['pkill', '-f', 'subscribe_workspaces'])
        subprocess.Popen(['pkill', '-f', 'subscribe_date'])
        subprocess.Popen(['pkill', '-f', 'subscribe_music'])
        subprocess.Popen(['pkill', '-f', 'subscribe_cpu'])
        subprocess.Popen(['pkill', '-f', 'subscribe_battery'])
        subprocess.Popen(['pkill', '-f', 'stalonetray'])
        subprocess.Popen(['killall', 'lemonbar'])

if __name__ == '__main__':
    atexit.register(Pylemon.kill_child_processes)
    instance = Pylemon()
