import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
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
    items: StorageProvisioner[] = [];
    @Output() createEvent = new EventEmitter();
    @Input() currentCluster: Cluster;

    ngOnInit(): void {
        this.refresh();
    }

    list() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    onCreate() {
        this.createEvent.emit();
    }

    refresh() {
        this.list();
    }
}
