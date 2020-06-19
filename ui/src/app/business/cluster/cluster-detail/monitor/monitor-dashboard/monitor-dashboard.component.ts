import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';

@Component({
    selector: 'app-monitor-dashboard',
    templateUrl: './monitor-dashboard.component.html',
    styleUrls: ['./monitor-dashboard.component.css']
})
export class MonitorDashboardComponent implements OnInit {

    constructor() {
    }

    @ViewChild('frame') frame: ElementRef;
    loading = true;


    ngOnInit(): void {
    }

    onFrameLoad() {
        this.frame.nativeElement.contentWindow.Mousetrap.unbindGlobal('esc');
        this.loading = false;
    }


}
