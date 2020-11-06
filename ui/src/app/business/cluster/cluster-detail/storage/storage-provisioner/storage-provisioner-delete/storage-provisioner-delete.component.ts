import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Cluster} from "../../../../cluster";
import {StorageProvisionerService} from "../storage-provisioner.service";
import {StorageProvisioner} from "../storage-provisioner";

@Component({
    selector: 'app-storage-provisioner-delete',
    templateUrl: './storage-provisioner-delete.component.html',
    styleUrls: ['./storage-provisioner-delete.component.css']
})
export class StorageProvisionerDeleteComponent implements OnInit {

    constructor(private service: StorageProvisionerService) {
    }

    opened = false;
    items: StorageProvisioner[] = [];
    @Output() deleted = new EventEmitter();
    @Input() currentCluster: Cluster;

    ngOnInit(): void {
    }

    open(items: StorageProvisioner[]) {
        this.opened = true;
        this.items = items;
    }

    onSubmit() {
        this.service.batch(this.currentCluster.name, this.items).subscribe(data => {
            this.opened = false;
            this.deleted.emit();
        });
    }

    onCancel() {
        this.opened = false;
    }

}
