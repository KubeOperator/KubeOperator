import {Component, EventEmitter, OnInit, ViewChild} from '@angular/core';
import {IstioService} from "./istio.service";
import {ActivatedRoute} from "@angular/router";
import {IstioHelper} from "./istios";
import {Cluster} from "../../cluster";
import {IstioListComponent} from "./istio-list/istio-list.component";
import {IstioDisableComponent} from "./istio-disable/istio-disable.component";

@Component({
    selector: 'app-istio',
    templateUrl: './istio.component.html',
    styleUrls: ['./istio.component.css']
})
export class IstioComponent implements OnInit {

    @ViewChild(IstioListComponent, {static: true})
    list: IstioListComponent;

    @ViewChild(IstioDisableComponent, {static: true})
    disable: IstioDisableComponent;

    constructor(private service: IstioService, private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openDisable(items: IstioHelper[]) {
        this.disable.open(items);
    }

    refresh() {
        this.list.refresh();
    }
}
