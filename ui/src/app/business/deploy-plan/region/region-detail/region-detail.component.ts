import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {BaseModelDirective} from '../../../../shared/class/BaseModelDirective';
import {Region} from '../region';
import {RegionService} from '../region.service';

@Component({
    selector: 'app-region-detail',
    templateUrl: './region-detail.component.html',
    styleUrls: ['./region-detail.component.css']
})
export class RegionDetailComponent extends BaseModelDirective<Region> implements OnInit {

    opened = false;
    @Output() detail = new EventEmitter();
    item: Region = new Region();

    constructor(private regionService: RegionService) {
        super(regionService)
    }

    ngOnInit(): void {

    }

    open(item) {
        this.opened = true;
        this.item = item;
    }

    cancel() {
        this.opened = false;
    }
}
