import {Component, Input, OnInit} from '@angular/core';
import {Registry} from "../registry";
import {RegistryService} from "../registry.service";
import {Cluster} from "../../../../cluster";

@Component({
    selector: 'app-registry-list',
    templateUrl: './registry-list.component.html',
    styleUrls: ['./registry-list.component.css']
})
export class RegistryListComponent implements OnInit {

    constructor(private registryService: RegistryService) {
    }

    loading = false;
    items: Registry[] = [];
    @Input() currentCluster: Cluster;

    ngOnInit(): void {
        this.refresh();
    }

    refresh() {
        this.loading = true;
        this.registryService.list(this.currentCluster.name).subscribe(data => {
            this.loading = false;
            this.items = [];
            for (const repository of data.repositories) {
                const item = new Registry();
                item.name = repository;
                this.items.push(item);
                this.loading = false;
            }
        });
    }

    loadChild(item: Registry) {
        item.loading = true;
        this.registryService.listTags(this.currentCluster.name, item.name).subscribe(data => {
            item.tags = data.tags;
            item.loading = false;
        });
    }


}
