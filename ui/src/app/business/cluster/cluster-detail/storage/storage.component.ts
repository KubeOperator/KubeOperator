import {Component, OnInit} from '@angular/core';
import {Cluster} from "../../cluster";
import {ClusterService} from "../../cluster.service";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-storage',
    templateUrl: './storage.component.html',
    styleUrls: ['./storage.component.css']
})
export class StorageComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

}
