import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Zone} from '../zone';
import {ZoneService} from '../zone.service';

@Component({
    selector: 'app-zone-list',
    templateUrl: './zone-list.component.html',
    styleUrls: ['./zone-list.component.css']
})
export class ZoneListComponent extends BaseModelComponent<Zone> implements OnInit {

    constructor(private zoneService: ZoneService) {
        super(zoneService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

}
