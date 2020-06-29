import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CreateStorageClassRequest} from '../../../storage';
import {NgForm} from '@angular/forms';

@Component({
    selector: 'app-storage-class-create-nfs',
    templateUrl: './storage-class-create-nfs.component.html',
    styleUrls: ['./storage-class-create-nfs.component.css']
})
export class StorageClassCreateNfsComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: CreateStorageClassRequest = new CreateStorageClassRequest();
    @ViewChild('nfsForm') nfsForm: NgForm;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.created.emit();
        this.opened = false;
    }

    open() {
        this.opened = true;
    }


}
