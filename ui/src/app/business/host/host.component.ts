import {Component, OnInit, ViewChild} from '@angular/core';
import {HostListComponent} from './host-list/host-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostDeleteComponent} from './host-delete/host-delete.component';
import {Host} from './host';

@Component({
    selector: 'app-host',
    templateUrl: './host.component.html',
    styleUrls: ['./host.component.css']
})
export class HostComponent implements OnInit {

    @ViewChild(HostListComponent, {static: true})
    list: HostListComponent;

    @ViewChild(HostCreateComponent, {static: true})
    create: HostCreateComponent;

    @ViewChild(HostDeleteComponent, {static: true})
    delete: HostDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: Host[]) {
        // this.delete.open(items);
    }

    openEdit(item: Host) {
        // this.edit.open(item);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }
}
