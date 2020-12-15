import {Component, OnInit, ViewChild} from '@angular/core';
import {IpPoolListComponent} from './ip-pool-list/ip-pool-list.component';
import {IpPoolCreateComponent} from './ip-pool-create/ip-pool-create.component';
import {IpPoolDeleteComponent} from './ip-pool-delete/ip-pool-delete.component';

@Component({
    selector: 'app-ip-pool',
    templateUrl: './ip-pool.component.html',
    styleUrls: ['./ip-pool.component.css']
})
export class IpPoolComponent implements OnInit {

    @ViewChild(IpPoolListComponent, {static: true})
    list: IpPoolListComponent;

    @ViewChild(IpPoolCreateComponent, {static: true})
    create: IpPoolCreateComponent;

    @ViewChild(IpPoolDeleteComponent, {static: true})
    delete: IpPoolDeleteComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }
}
