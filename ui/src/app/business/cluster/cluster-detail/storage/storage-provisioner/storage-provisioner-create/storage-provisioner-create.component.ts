import {Component, OnInit} from '@angular/core';
import {CreateStorageProvisionerRequest} from '../storage-provisioner';

@Component({
    selector: 'app-storage-provisioner-create',
    templateUrl: './storage-provisioner-create.component.html',
    styleUrls: ['./storage-provisioner-create.component.css']
})
export class StorageProvisionerCreateComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item = new CreateStorageProvisionerRequest();

    ngOnInit(): void {
    }


}
