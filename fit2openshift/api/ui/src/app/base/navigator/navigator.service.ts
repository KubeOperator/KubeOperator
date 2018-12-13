import { Injectable } from '@angular/core';
import { Subject } from "rxjs/index";
import { Operation } from "./navigator";

@Injectable({
  providedIn: 'root'
})
export class NavigatorService {
  private message$ = new Subject<Operation>();
  messageObserver = this.message$.asObservable();

  constructor() { }

  showProjectsNav() {
    const operation: Operation = {
      action: "show",
      target: "projects"
    };
    this.message$.next(operation)
  }

  showDetailNav(projectId: string) {
    const operation: Operation = {
      action: "show",
      target: "detail",
      meta: projectId
    };
    this.message$.next(operation)
  }
}
