import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';

@Component({
    selector: 'app-zone-detail',
    templateUrl: './zone-detail.component.html',
    styleUrls: ['./zone-detail.component.css']
})
export class ZoneDetailComponent extends BaseModelComponent<Zone> implements OnInit {

    opened = false;
    item: Zone = new Zone();

    @Output() detail = new EventEmitter();

    constructor(private zoneService: ZoneService) {
        super(zoneService);
    }

    ngOnInit(): void {
    }


    open(item) {
        this.item = item;
        this.opened = true;
    }

    cancel() {
        this.opened = false;
    }
}
