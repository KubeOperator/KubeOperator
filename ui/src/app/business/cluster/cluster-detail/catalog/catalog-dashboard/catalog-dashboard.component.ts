import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {Cluster} from "../../../cluster";
import {ClusterTool} from "../../tools/tools";
import {ActivatedRoute} from "@angular/router";
import {ClusterService} from "../../../cluster.service";
import {DomSanitizer} from "@angular/platform-browser";

@Component({
    selector: 'app-catelog-dashboard',
    templateUrl: './catalog-dashboard.component.html',
    styleUrls: ['./catalog-dashboard.component.css']
})
export class CatalogDashboardComponent implements OnInit {

    @ViewChild('frame') frame: ElementRef;
    loading = true;
    ready = false;
    @Input() currentCluster: Cluster;
    @Input() item: ClusterTool;
    url: any;

    constructor(private route: ActivatedRoute, private clusterService: ClusterService, private sanitizer: DomSanitizer) {
    }

    ngOnInit(): void {
        this.refresh();
    }
    refresh() {
        this.clusterService.secret(this.currentCluster.name).subscribe(data => {
            localStorage.setItem('kubeapps_auth_token', data.kubernetesToken);
            localStorage.setItem('kubeapps_auth_token_oid', 'false');
            const url = `/proxy/kubeapps/${!this.currentCluster.name}/root#/ns/default/apps`;
            this.url = this.sanitizer.bypassSecurityTrustResourceUrl(url);
            this.ready = true;
        });
    }
    onFrameLoad() {
        this.loading = false;
    }

}
