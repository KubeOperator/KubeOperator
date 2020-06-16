import {Component, OnInit} from '@angular/core';
import {BaseModelComponent} from '../../../../../shared/class/BaseModelComponent';
import {Task} from '../task';
import {TaskService} from '../task.service';
import {Cluster} from '../../../cluster';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-task-list',
    templateUrl: './task-list.component.html',
    styleUrls: ['./task-list.component.css']
})
export class TaskListComponent extends BaseModelComponent<Task> implements OnInit {

    currentCluster: Cluster;

    constructor(service: TaskService, private route: ActivatedRoute) {
        super(service);
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster =data.cluster;
            this.service.variable.set('cluster_name', this.currentCluster.name);
        });
    }

}
