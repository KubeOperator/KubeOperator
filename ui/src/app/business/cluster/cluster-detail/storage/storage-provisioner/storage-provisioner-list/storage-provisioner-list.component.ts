import {Component, Input, OnInit} from '@angular/core';
import {StorageProvisionerService} from '../storage-provisioner.service';
import {Cluster} from '../../../../cluster';
import {StorageProvisioner} from '../storage-provisioner';

@Component({
    selector: 'app-storage-provisioner-list',
    templateUrl: './storage-provisioner-list.component.html',
    styleUrls: ['./storage-provisioner-list.component.css']
})
export class StorageProvisionerListComponent implements OnInit {

    constructor(private service: StorageProvisionerService) {
    }

    loading = false;
    @Input() currentCluster: Cluster;
    items: StorageProvisioner[] = [];

    ngOnInit(): void {
        this.refresh();
    }

    list() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    refresh() {
        this.list();
    }
}
