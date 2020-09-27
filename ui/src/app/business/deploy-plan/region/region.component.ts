import {Component, OnInit, ViewChild} from '@angular/core';
import {RegionListComponent} from './region-list/region-list.component';
import {RegionCreateComponent} from './region-create/region-create.component';
import {RegionDeleteComponent} from './region-delete/region-delete.component';
import {BaseModelDirective} from '../../../shared/class/BaseModelDirective';
import {Region} from './region';
import {RegionService} from './region.service';
import {RegionDetailComponent} from './region-detail/region-detail.component';

@Component({
    selector: 'app-region',
    templateUrl: './region.component.html',
    styleUrls: ['./region.component.css']
})
export class RegionComponent extends BaseModelDirective<Region> implements OnInit {

    @ViewChild(RegionListComponent, {static: true})
    list: RegionListComponent;

    @ViewChild(RegionCreateComponent, {static: true})
    create: RegionCreateComponent;

    @ViewChild(RegionDeleteComponent, {static: true})
    delete: RegionDeleteComponent;

    @ViewChild(RegionDetailComponent, {static: true})
    detail: RegionDetailComponent;

    constructor(private regionService: RegionService) {
        super(regionService);
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
}
