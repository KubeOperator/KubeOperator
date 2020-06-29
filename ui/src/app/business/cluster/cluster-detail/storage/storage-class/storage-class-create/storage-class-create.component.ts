import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CreateStorageClassRequest} from '../../storage';
import {NgForm} from '@angular/forms';

@Component({
    selector: 'app-storage-class-create',
    templateUrl: './storage-class-create.component.html',
    styleUrls: ['./storage-class-create.component.css']
})
export class StorageClassCreateComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item = new CreateStorageClassRequest();
    @ViewChild('scForm') scForm: NgForm;
    @Output() selected = new EventEmitter<CreateStorageClassRequest>();

    ngOnInit(): void {
    }

    reset() {
        this.item = new CreateStorageClassRequest();
        this.scForm.resetForm();
    }

    open() {
        this.opened = true;
        this.reset();
    }

    onSubmit() {
        this.opened = false;
        this.selected.emit(this.item);
    }

    onCancel() {
        this.opened = false;
    }


}
