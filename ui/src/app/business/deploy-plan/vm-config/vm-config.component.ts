import {Component, OnInit, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {VmConfig} from './vm-config';
import {VmConfigService} from './vm-config.service';
import {VmConfigListComponent} from './vm-config-list/vm-config-list.component';
import {VmConfigDeleteComponent} from './vm-config-delete/vm-config-delete.component';
import {VmConfigCreateComponent} from './vm-config-create/vm-config-create.component';
import {VmConfigUpdateComponent} from './vm-config-update/vm-config-update.component';

@Component({
    selector: 'app-vm-config',
    templateUrl: './vm-config.component.html',
    styleUrls: ['./vm-config.component.css']
})
export class VmConfigComponent extends BaseModelDirective<VmConfig> implements OnInit {

    @ViewChild(VmConfigListComponent, {static: true})
    list: VmConfigListComponent;

    @ViewChild(VmConfigDeleteComponent, {static: true})
    delete: VmConfigDeleteComponent;

    @ViewChild(VmConfigCreateComponent, {static: true})
    create: VmConfigCreateComponent;

    @ViewChild(VmConfigUpdateComponent, {static: true})
    update: VmConfigUpdateComponent;

    constructor(private vmConfigService: VmConfigService) {
        super(vmConfigService);
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
