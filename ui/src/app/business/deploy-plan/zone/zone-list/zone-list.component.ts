import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';
import {Region} from '../../region/region';

@Component({
    selector: 'app-zone-list',
    templateUrl: './zone-list.component.html',
    styleUrls: ['./zone-list.component.css']
})
export class ZoneListComponent extends BaseModelDirective<Zone> implements OnInit {

    @Output() detailEvent = new EventEmitter<Region>();

    constructor(private zoneService: ZoneService) {
        super(zoneService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onDetail(item) {
        this.detailEvent.emit(item);
    }
}
