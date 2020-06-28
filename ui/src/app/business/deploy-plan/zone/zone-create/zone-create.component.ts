import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Zone, ZoneCreateRequest} from '../zone';
import {ZoneService} from '../zone.service';
import {RegionService} from '../../region/region.service';
import {Region} from '../../region/region';

@Component({
    selector: 'app-zone-create',
    templateUrl: './zone-create.component.html',
    styleUrls: ['./zone-create.component.css']
})
export class ZoneCreateComponent extends BaseModelComponent<Zone> implements OnInit {

    opened = false;
    item: ZoneCreateRequest = new ZoneCreateRequest();
    regions: Region[] = [];

    constructor(private zoneService: ZoneService, private regionService: RegionService) {
        super(zoneService);
    }

    ngOnInit(): void {

    }

    open() {
        this.item = new ZoneCreateRequest();
        this.opened = true;
    }

    onCancel(): void {
        this.opened = false;
    }

    onSubmit(): void {

    }

    listRegion() {
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }

}
