import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {KubernetesService} from "../../../kubernetes.service";
import {ClusterTool} from "../tools";
import {V1StorageClass} from "@kubernetes/client-node";
import {NgForm} from "@angular/forms";
import {Cluster} from "../../../cluster";
import {ToolsService} from "../tools.service";

@Component({
    selector: 'app-tools-enable',
    templateUrl: './tools-enable.component.html',
    styleUrls: ['./tools-enable.component.css']
})
export class ToolsEnableComponent implements OnInit {

    constructor(private kubernetesService: KubernetesService, private toolsService: ToolsService) {
    }

    opened = false;
    isSubmitGoing = false;
    item: ClusterTool = new ClusterTool();
    storageClazz: V1StorageClass[] = [];
    @ViewChild('itemForm') itemForm: NgForm;
    @Output() enabled = new EventEmitter();
    @Input() currentCluster: Cluster;


    ngOnInit(): void {
    }

    onSubmit() {
        this.toolsService.enable(this.currentCluster.name, this.item).subscribe(data => {
            this.opened = false;
            this.enabled.emit();
        });
    }

    onCancel() {
        this.opened = false;
    }

    reset() {
        this.itemForm.resetForm();
        this.listStorageClass();
    }

    open(item: ClusterTool) {
        this.reset();
        this.opened = true;
        this.setDefaultVars(item);
        this.item = item;
        console.log(this.item);
    }

    listStorageClass() {
        this.kubernetesService.listStorageClass(this.currentCluster.name, '', true).subscribe(data => {
            this.storageClazz = data.items;
        });
    }

    setDefaultVars(item: ClusterTool) {
        switch (item.name) {
            case 'prometheus':
                item.vars = {
                    'server.retention': 10,
                    'server.persistentVolume.enabled': false,
                    'server.persistentVolume.size': 10,
                    'server.persistentVolume.storageClass': '',
                };
                break;
            case 'chartmuseum':
                item.vars = {
                    'persistence.enabled': false,
                    'env.open.DISABLE_API': false,
                    'persistence.storageClass': '',
                    'persistence.size': 10,
                };
                break;
            case 'registry':
                item.vars = {
                    'persistence.enabled': false,
                    'persistence.storageClass': '',
                    'service.type': 'NodePort',
                    'persistence.size': 10,
                };
                break;
            case 'efk':
                item.vars = {
                    'elasticsearch.persistence.enabled': false,
                    'elasticsearch.volumeClaimTemplate.resources.requests.storage': 10,
                    'elasticsearch.volumeClaimTemplate.storageClassName': '',
                };
                break;
            case 'kubeapps':
                item.vars = {
                    'postgresql.persistence.enabled': false,
                    'postgresql.persistence.size': 10,
                    'global.storageClass': ''
                };
                break;
            case 'dashboard':
                item.vars = {};
                break;
        }
    }

}
