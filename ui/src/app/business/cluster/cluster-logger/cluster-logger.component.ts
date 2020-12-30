import {Component, ElementRef, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {Terminal} from "xterm";
import {ClusterLoggerService} from "./cluster-logger.service";

@Component({
    selector: 'app-cluster-logger',
    templateUrl: './cluster-logger.component.html',
    styleUrls: ['./cluster-logger.component.css']
})
export class ClusterLoggerComponent implements OnInit, OnDestroy {

    term: Terminal;
    timer;
    isRun: boolean = true; 
    @ViewChild('terminal', {static: true}) terminal: ElementRef;


    constructor(private loggerService: ClusterLoggerService) {
    }

    ngOnDestroy() {
        if (this.timer) {
            clearInterval(this.timer);
        }
    }

    ngOnInit() {
        this.term = new Terminal({
            disableStdin: true,
            cursorStyle: 'bar',
            cols: 100,
            rows: 58,
            letterSpacing: 1,
            fontSize: 12
        });

        this.term.open(this.terminal.nativeElement);
        this.term.write('connect to logger...');
        setTimeout(() => {
            this.term.clear();
        }, 3000);
        const clusterName = this.getQueryVariable('clusterName');
        const logId = this.getQueryVariable('logId');
        this.timer = setInterval(() => {
            if (this.isRun) {
                if (logId) {
                    this.loggerService.getProvisionerLog(clusterName, logId).subscribe(data => {
                        this.term.clear();
                        const text = data.msg.replace(/\n/g, '\r\n');
                        this.term.write(text);
                        setTimeout(() => {
                            this.term.scrollToBottom();
                        }, 100);
                    }, error => {
                        this.term.write('no log to show');
                    });
                } else {
                    this.loggerService.getClusterLog(clusterName).subscribe(data => {
                        this.term.clear();
                        const text = data.msg.replace(/\n/g, '\r\n');
                        this.term.write(text);
                        setTimeout(() => {
                            this.term.scrollToBottom();
                        }, 100);
                    }, error => {
                        this.term.write('no log to show');
                    });
                }
            }
        }, 5000);
    }

    changeMode() {
        this.isRun = !this.isRun;
    }

    getQueryVariable(variable): string {
        const query = window.location.search.substring(1);
        const vars = query.split('&');
        for (const v of vars) {
            const pair = v.split('=');
            if (pair[0] === variable) {
                return pair[1];
            }
        }
        return null;
    }
}
