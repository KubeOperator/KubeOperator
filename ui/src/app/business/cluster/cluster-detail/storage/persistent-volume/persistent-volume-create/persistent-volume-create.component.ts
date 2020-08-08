import {Component, EventEmitter, OnInit, Output} from '@angular/core';

@Component({
    selector: 'app-persistent-volume-create',
    templateUrl: './persistent-volume-create.component.html',
    styleUrls: ['./persistent-volume-create.component.css']
})
export class PersistentVolumeCreateComponent implements OnInit {

    constructor() {
    }

    opened = false;
    provisioner = '';

    @Output() selected = new EventEmitter();

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.selected.emit(this.provisioner);
        this.opened = false;
    }
}
