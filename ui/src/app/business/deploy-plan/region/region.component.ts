import {Component, OnInit, ViewChild} from '@angular/core';
import {RegionListComponent} from "./region-list/region-list.component";
import {RegionCreateComponent} from "./region-create/region-create.component";
import {RegionDeleteComponent} from "./region-delete/region-delete.component";

@Component({
    selector: 'app-region',
    templateUrl: './region.component.html',
    styleUrls: ['./region.component.css']
})
export class RegionComponent implements OnInit {

    @ViewChild(RegionListComponent, {static: true})
    list: RegionListComponent;

    @ViewChild(RegionCreateComponent, {static: true})
    create: RegionCreateComponent;

    @ViewChild(RegionDeleteComponent, {static: true})
    delete: RegionDeleteComponent;


    constructor() {
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
}
