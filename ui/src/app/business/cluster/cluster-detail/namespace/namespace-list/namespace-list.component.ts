import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {V1Namespace} from '@kubernetes/client-node';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/app-alert/alert';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';

@Component({
    selector: 'app-namespace-list',
    templateUrl: './namespace-list.component.html',
    styleUrls: ['./namespace-list.component.css']
})
export class NamespaceListComponent implements OnInit {

    loading = true;
    selected = [];
    items: V1Namespace[] = [];
    page = 1;
    currentCluster: Cluster;
    @Output() deleteEvent = new EventEmitter<string>();
    @Output() createEvent = new EventEmitter<string>();

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listNamespaces(this.currentCluster.name).subscribe(data => {
            this.loading = false;
            this.items = data.items;
        });
    }

    onDelete(item: V1Namespace) {
        this.deleteEvent.emit(item.metadata.name);
    }

    onCreate() {
        this.createEvent.emit();
    }

}
