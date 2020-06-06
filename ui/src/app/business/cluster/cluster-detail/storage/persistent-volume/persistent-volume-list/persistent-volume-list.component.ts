import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {Cluster} from "../../../../cluster";
import {V1Namespace, V1PersistentVolume} from "@kubernetes/client-node";
import {KubernetesService} from "../../../../kubernetes.service";

@Component({
    selector: 'app-persistent-volume-list',
    templateUrl: './persistent-volume-list.component.html',
    styleUrls: ['./persistent-volume-list.component.css']
})
export class PersistentVolumeListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1PersistentVolume[] = [];
    loading = true;
    selected = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listPersistentVolumes(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            console.log(data);
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

}
