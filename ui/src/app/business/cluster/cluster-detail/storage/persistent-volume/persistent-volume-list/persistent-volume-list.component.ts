import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1PersistentVolume} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';

@Component({
    selector: 'app-persistent-volume-list',
    templateUrl: './persistent-volume-list.component.html',
    styleUrls: ['./persistent-volume-list.component.css']
})
export class PersistentVolumeListComponent implements OnInit {

    items: V1PersistentVolume[] = [];
    loading = true;
    selected = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';
    @Output() createEvent = new EventEmitter();
    @Input() currentCluster: Cluster;
    @Output() deleteEvent = new EventEmitter<string>();

    constructor(private service: KubernetesService) {
    }

    ngOnInit(): void {
        this.list();
    }

    list() {
        this.loading = true;
        this.service.listPersistentVolumes(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }


    getSource(item: V1PersistentVolume) {
        return item.spec.storageClassName;
    }

    refresh() {
        this.list();
    }

    onCreate() {
        this.createEvent.emit();
    }

    onDelete(name) {
        this.deleteEvent.emit(name);
    }
}
