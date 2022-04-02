import {Component, OnInit} from '@angular/core';
import {KubernetesService} from '../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';
import {EventService} from './event.service';
import {CommonAlertService} from '../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../layout/common-alert/alert';
import {ClrDatagridSortOrder} from '@clr/angular';

@Component({
    selector: 'app-event',
    templateUrl: './event.component.html',
    styleUrls: ['./event.component.css']
})
export class EventComponent implements OnInit {

    loading = false;
    currentCluster: Cluster;
    namespaces;
    events;
    currentNamespace: string;
    npdExists = false;
    nextToken = '';
    previousToken = '';
    continueToken = '';
    showPage = true;
    descSort = ClrDatagridSortOrder.ASC;

    constructor(private kubernetesService: KubernetesService,
                private route: ActivatedRoute,
                private eventService: EventService,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
        this.loading = true;
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            let search = {
                kind: "namespacelist",
                cluster: this.currentCluster.name,
                continue: "",
                limit: 0,
                namespace: "",
                name: "",
            }
            this.kubernetesService.listResource(search).subscribe(res => {
                this.namespaces = res.items;
                if (this.namespaces.length > 0) {
                    const namespace = this.namespaces[0];
                    this.currentNamespace = namespace.metadata.name;
                    this.listEvents(this.currentNamespace);
                }
            });
            this.getNpdExists();
        });
    }

    changeNamespace(namespace) {
        this.currentNamespace = namespace;
        this.nextToken = '';
        this.previousToken = '';
        this.continueToken = '';
        this.showPage = false;
        setTimeout(x => this.showPage = true);
        this.listEvents(namespace);
    }

    listEvents(namespace) {
        this.loading = true;
        let search = {
            kind: "eventlist",
            cluster: this.currentCluster.name,
            continue: this.continueToken,
            limit: 30,
            namespace: namespace,
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(res => {
            this.events = res.items;
            this.loading = false;
            this.nextToken = res.metadata[this.kubernetesService.continueTokenKey] ? res.metadata[this.kubernetesService.continueTokenKey] : '';
        });
    }

    getNpdExists() {
        let search = {
            kind: "podlist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.kubernetesService.listResource(search).subscribe(data => {
            const pods = data.items;
            for (const pod of pods) {
                if (pod.metadata.generateName === 'node-problem-detector-') {
                    this.npdExists = true;
                    break;
                }
            }
        });
    }

    changeNpd(exists) {
        this.npdExists = !exists;
        let op = 'create';
        if (exists) {
            op = 'delete';
        }
        this.eventService.changeNpd(this.currentCluster.name, op).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_UPDATE_SUCCESS'), AlertLevels.SUCCESS);
        }, error => {
            this.npdExists = exists;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
