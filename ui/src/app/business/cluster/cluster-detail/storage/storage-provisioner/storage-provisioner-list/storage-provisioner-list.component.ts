import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {StorageProvisionerService} from '../storage-provisioner.service';
import {Cluster} from '../../../../cluster';
import {StorageProvisioner} from '../storage-provisioner';
import {ClusterLoggerService} from "../../../../cluster-logger/cluster-logger.service";

@Component({
    selector: 'app-storage-provisioner-list',
    templateUrl: './storage-provisioner-list.component.html',
    styleUrls: ['./storage-provisioner-list.component.css']
})
export class StorageProvisionerListComponent implements OnInit {

    timer;
    constructor(private service: StorageProvisionerService, private loggerService: ClusterLoggerService) {
    }

    loading = false;
    items: StorageProvisioner[] = [];
    selected: StorageProvisioner[] = [];
    opened = false;
    detailItem: StorageProvisioner = new StorageProvisioner();
    @Output() createEvent = new EventEmitter();
    @Output() deleteEvent = new EventEmitter();
    @Input() currentCluster: Cluster;

    ngOnInit(): void {
        this.refresh();
        this.polling();
    }

    ngOnDestroy(): void {
        clearInterval(this.timer);
    }

    list() {
        this.service.list(this.currentCluster.name).subscribe(data => {
            this.items = data;
        });
    }

    onCreate() {
        this.createEvent.emit();
    }

    onDelete() {
        this.deleteEvent.emit(this.selected);
    }

    refresh() {
        this.list();
    }

    onShowLogger(item: StorageProvisioner) {
        this.loggerService.openProvisionerLogger(this.currentCluster.name, item.id);
    }

    polling() {
        this.timer = setInterval(() => {
            this.refresh();
        }, 5000);
    }

    openMessage(item) {
        this.opened = true;
        this.detailItem = item;
    }
}
