import {Component, OnInit} from '@angular/core';
import {KubernetesService} from "../../../kubernetes.service";
import {V1Namespace, V1Node} from "@kubernetes/client-node";
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-node-list',
    templateUrl: './node-list.component.html',
    styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

    loading = true;
    selected = [];
    items: V1Node[] = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';
    page = 1;
    currentCluster: Cluster;

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listNodes(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

}
