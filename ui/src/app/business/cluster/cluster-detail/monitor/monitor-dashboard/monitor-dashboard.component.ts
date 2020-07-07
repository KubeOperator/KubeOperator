import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster, ClusterMonitor} from '../../../cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../../../cluster.service';
import {DomSanitizer} from '@angular/platform-browser';

@Component({
    selector: 'app-monitor-dashboard',
    templateUrl: './monitor-dashboard.component.html',
    styleUrls: ['./monitor-dashboard.component.css']
})
export class MonitorDashboardComponent implements OnInit {


    @ViewChild('frame') frame: ElementRef;
    loading = true;
    @Input() currentCluster: Cluster;
    monitor: ClusterMonitor;
    url: any;

    constructor(private route: ActivatedRoute, private clusterService: ClusterService, private sanitizer: DomSanitizer) {
    }

    ngOnInit(): void {
    }

    refresh() {
        this.clusterService.monitor(this.currentCluster.name).subscribe(data => {
            const url = data.dashboardUrl + '?orgId=1&kiosk';
            this.url = this.sanitizer.bypassSecurityTrustResourceUrl(url);
            console.log(this.url);
            this.monitor = data;
        });
    }

    onFrameLoad() {
        this.frame.nativeElement.contentWindow.Mousetrap.unbindGlobal('esc');
        this.loading = false;
    }


}
