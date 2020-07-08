import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelComponent} from '../../../../shared/class/BaseModelComponent';
import {Region} from '../region';
import {RegionService} from '../region.service';

@Component({
    selector: 'app-region-list',
    templateUrl: './region-list.component.html',
    styleUrls: ['./region-list.component.css']
})
export class RegionListComponent extends BaseModelComponent<Region> implements OnInit {

    @Output() detailEvent = new EventEmitter<Region>();


    constructor(private regionService: RegionService) {
        super(regionService);
    }

    ngOnInit(): void {
        super.ngOnInit();
    }

    onDetail(item) {
        this.detailEvent.emit(item);
    }
}
