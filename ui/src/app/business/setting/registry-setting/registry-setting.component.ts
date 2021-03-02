import {Component, OnInit, ViewChild} from '@angular/core';
import {RegistryListComponent} from './registry-list/registry-list.component';
import {RegistryCreateComponent} from './registry-create/registry-create.component';
import {RegistryDeleteComponent} from './registry-delete/registry-delete.component';
import {RegistryUpdateComponent} from './registry-update/registry-update.component';

@Component({
    selector: 'app-registry-setting',
    templateUrl: './registry-setting.component.html',
    styleUrls: ['./registry-setting.component.css']
})
export class RegistrySettingComponent implements OnInit {

    @ViewChild(RegistryListComponent, {static: true})
    list: RegistryListComponent;

    @ViewChild(RegistryCreateComponent, {static: true})
    create: RegistryCreateComponent;

    @ViewChild(RegistryDeleteComponent, {static: true})
    delete: RegistryDeleteComponent;

    @ViewChild(RegistryUpdateComponent, {static: true})
    update: RegistryUpdateComponent;

    constructor() {
    }

    ngOnInit(): void {

    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items) {
        this.delete.open(items);
    }

    openUpdate(item) {
        this.update.open(item);
    }
}
