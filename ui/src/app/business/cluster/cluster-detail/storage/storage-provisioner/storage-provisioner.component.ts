import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';

@Component({
    selector: 'app-storage-provisioner',
    templateUrl: './storage-provisioner.component.html',
    styleUrls: ['./storage-provisioner.component.css']
})
export class StorageProvisionerComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

}
