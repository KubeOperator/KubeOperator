import {Component, OnInit, ViewChild} from '@angular/core';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Zone} from './zone';
import {ZoneService} from './zone.service';
import {ZoneListComponent} from './zone-list/zone-list.component';
import {ZoneUpdateComponent} from './zone-update/zone-update.component';
import {ZoneDeleteComponent} from './zone-delete/zone-delete.component';
import {ZoneCreateComponent} from './zone-create/zone-create.component';
import {ZoneDetailComponent} from './zone-detail/zone-detail.component';

@Component({
    selector: 'app-zone',
    templateUrl: './zone.component.html',
    styleUrls: ['./zone.component.css']
})
export class ZoneComponent extends BaseModelDirective<Zone> implements OnInit {

    @ViewChild(ZoneListComponent, {static: true})
    list: ZoneListComponent;

    @ViewChild(ZoneUpdateComponent, {static: true})
    update: ZoneUpdateComponent;

    @ViewChild(ZoneDeleteComponent, {static: true})
    delete: ZoneDeleteComponent;

    @ViewChild(ZoneCreateComponent, {static: true})
    create: ZoneCreateComponent;

    @ViewChild(ZoneDetailComponent, {static: true})
    detail: ZoneDetailComponent;

    constructor(private zoneService: ZoneService) {
        super(zoneService);
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

    openDetail(item) {
        this.detail.open(item);
    }

    openUpdate(item) {
        this.update.open(item);
    }
}
