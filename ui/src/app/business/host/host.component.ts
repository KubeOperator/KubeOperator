import {Component, OnInit, ViewChild} from '@angular/core';
import {HostListComponent} from './host-list/host-list.component';
import {HostCreateComponent} from './host-create/host-create.component';
import {HostDeleteComponent} from './host-delete/host-delete.component';
import {Host} from './host';
import {HostDetailComponent} from './host-detail/host-detail.component';
import {HostStatusDetailComponent} from './host-status-detail/host-status-detail.component';
import {HostImportComponent} from './host-import/host-import.component';

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

    @ViewChild(HostDetailComponent, {static: true})
    detail: HostDetailComponent;

    @ViewChild(HostStatusDetailComponent, {static: true})
    statusDetail: HostStatusDetailComponent;

    @ViewChild(HostImportComponent, {static: true})
    import: HostImportComponent;

    constructor() {
    }

    ngOnInit(): void {
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: Host[]) {
        this.delete.open(items);
    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openDetail(item) {
        this.detail.open(item);
    }

    openStatusDetail(item: Host) {
        this.statusDetail.open(item);
    }

    openImport() {
        this.import.open();
    }
}
