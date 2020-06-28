import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {CloudZoneRequest, Zone, ZoneCreateRequest} from '../zone';
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
    cloudZoneRequest: CloudZoneRequest = new CloudZoneRequest();
    regions: Region[] = [];
    cloudZones: [] = [];
    region: Region = new Region();
    @Output() created = new EventEmitter();


    constructor(private zoneService: ZoneService, private regionService: RegionService) {
        super(zoneService);
    }

    ngOnInit(): void {

    }

    open() {
        this.item = new ZoneCreateRequest();
        this.opened = true;
        this.listRegion();
    }

    onCancel(): void {
        this.opened = false;
    }

    onSubmit(): void {

    }

    changeRegion() {
        this.regions.forEach(region => {
            if (region.name === this.item.region) {
                this.region = region;
                this.region.regionVars = JSON.parse(this.region.vars);
                this.cloudZoneRequest.cloudVars = JSON.parse(this.region.vars);
            }
        });
    }

    listRegion() {
        this.regionService.list().subscribe(res => {
            this.regions = res.items;
        }, error => {

        });
    }

    onBasicFormCommit(){
        this.loading = true;
        this.cloudZoneRequest.datacenter = this.region.datacenter;
        this.zoneService.listClusters(this.cloudZoneRequest).subscribe(res => {
            this.cloudZones = res.result;
            this.loading = false;
        })
    }

}
