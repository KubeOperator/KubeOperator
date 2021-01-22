import {Component, OnInit, ViewChild} from '@angular/core';
import {PersistentVolumeListComponent} from './persistent-volume-list/persistent-volume-list.component';
import {PersistentVolumeCreateComponent} from './persistent-volume-create/persistent-volume-create.component';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {PersistentVolumeCreateHostPathComponent} from './persistent-volume-create/persistent-volume-create-host-path/persistent-volume-create-host-path.component';
import {PersistentVolumeCreateLocalStorageComponent} from "./persistent-volume-create/persistent-volume-create-local-storage/persistent-volume-create-local-storage.component";
import {PersistentVolumeDeleteComponent} from './persistent-volume-delete/persistent-volume-delete.component';

@Component({
    selector: 'app-persistent-volume',
    templateUrl: './persistent-volume.component.html',
    styleUrls: ['./persistent-volume.component.css']
})
export class PersistentVolumeComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    currentCluster: Cluster;

    @ViewChild(PersistentVolumeListComponent, {static: true})
    list: PersistentVolumeListComponent;

    @ViewChild(PersistentVolumeCreateComponent, {static: true})
    create: PersistentVolumeCreateComponent;

    @ViewChild(PersistentVolumeDeleteComponent, {static: true})
    delete: PersistentVolumeDeleteComponent;

    @ViewChild(PersistentVolumeCreateHostPathComponent, {static: true})
    hostPath: PersistentVolumeCreateHostPathComponent;

    @ViewChild(PersistentVolumeCreateLocalStorageComponent, {static: true})
    localStorage: PersistentVolumeCreateLocalStorageComponent;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openCreate() {
        this.create.open();
    }

    openSelected(provisioner: string) {
        switch (provisioner) {
            case 'localStorage':
                this.localStorage.open();
                break;
            case 'hostPath':
                this.hostPath.open();
                break;
        }
    }

    refresh() {
        this.list.refresh();
    }


    openDelete(name) {
        this.delete.open(name);
    }
}
