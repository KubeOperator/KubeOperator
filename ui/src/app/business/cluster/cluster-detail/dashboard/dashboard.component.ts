import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster, ClusterMonitor} from "../../cluster";
import {ActivatedRoute} from "@angular/router";
import {ClusterService} from "../../cluster.service";
import {DomSanitizer} from "@angular/platform-browser";

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {

    @ViewChild('frame') frame: ElementRef;
    loading = true;
    @Input() currentCluster: Cluster;
    url: any;

    constructor(private route: ActivatedRoute, private sanitizer: DomSanitizer) {
    }

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.url = this.sanitizer.bypassSecurityTrustResourceUrl('/proxy/dashboard/test/root');
        console.log(this.url);
    }

    onFrameLoad() {
        this.frame.nativeElement.contentWindow.Mousetrap.unbindGlobal('esc');
        this.loading = false;
    }

}
