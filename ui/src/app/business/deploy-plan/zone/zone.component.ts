import {Component, OnInit, ViewChild} from '@angular/core';
import {BaseModelComponent} from '../../../shared/class/BaseModelComponent';
import {Zone} from './zone';
import {ZoneService} from './zone.service';
import {ZoneListComponent} from './zone-list/zone-list.component';
import {ZoneUpdateComponent} from './zone-update/zone-update.component';
import {ZoneDeleteComponent} from './zone-delete/zone-delete.component';
import {ZoneCreateComponent} from './zone-create/zone-create.component';

@Component({
    selector: 'app-zone',
    templateUrl: './zone.component.html',
    styleUrls: ['./zone.component.css']
})
export class ZoneComponent extends BaseModelComponent<Zone> implements OnInit {

    @ViewChild(ZoneListComponent, {static: true})
    list: ZoneListComponent;

    @ViewChild(ZoneUpdateComponent, {static: true})
    update: ZoneUpdateComponent;

    @ViewChild(ZoneDeleteComponent, {static: true})
    delete: ZoneDeleteComponent;

    @ViewChild(ZoneCreateComponent, {static: true})
    create: ZoneCreateComponent;


    constructor(private zoneService: ZoneService) {
        super(zoneService);
    }

    ngOnInit(): void {

    }

    refresh() {
        this.list.reset();
        this.list.refresh();
    }
}
