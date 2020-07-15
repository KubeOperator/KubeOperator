import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from '../../../cluster';
import {ActivatedRoute} from '@angular/router';
import {ClusterService} from '../../../cluster.service';
import {DomSanitizer} from '@angular/platform-browser';
import {ClusterTool} from "../../tools/tools";

@Component({
    selector: 'app-monitor-dashboard',
    templateUrl: './monitor-dashboard.component.html',
    styleUrls: ['./monitor-dashboard.component.css']
})
export class MonitorDashboardComponent implements OnInit {


    @ViewChild('frame') frame: ElementRef;
    loading = true;
    @Input() currentCluster: Cluster;
    @Input() item: ClusterTool;
    url: any;

    constructor(private route: ActivatedRoute, private clusterService: ClusterService, private sanitizer: DomSanitizer) {
    }

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        const url = this.item.vars['url'] + '?orgId=1&kiosk';
        console.log(this.item);
        this.url = this.sanitizer.bypassSecurityTrustResourceUrl(url);
    }

    onFrameLoad() {
        this.frame.nativeElement.contentWindow.Mousetrap.unbindGlobal('esc');
        this.loading = false;
    }


}
