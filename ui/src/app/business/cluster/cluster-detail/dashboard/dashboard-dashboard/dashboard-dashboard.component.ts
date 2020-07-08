import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";
import {DomSanitizer} from "@angular/platform-browser";

@Component({
    selector: 'app-dashboard-dashboard',
    templateUrl: './dashboard-dashboard.component.html',
    styleUrls: ['./dashboard-dashboard.component.css']
})
export class DashboardDashboardComponent implements OnInit {

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
    }

    onFrameLoad() {
        this.loading = false;
    }

}
