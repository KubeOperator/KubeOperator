import { Component, OnInit, Input, Output, ElementRef, EventEmitter, ViewChild, AfterViewInit, OnDestroy } from '@angular/core';
import { Terminal } from 'xterm';

@Component({
  selector: 'app-celery-log',
  templateUrl: './celery-log.component.html',
  styles: []
})
export class CeleryLogComponent implements OnInit, AfterViewInit, OnDestroy {
  @Input() show: boolean;
  @Input() url: string;
  @Output() closeModal = new EventEmitter<boolean>();
  @ViewChild('terminal')
  terminal: ElementRef;
  closable = false;
  size = 'lg';
  staticBackdrop = true;
  ws: WebSocket;
  term: Terminal;

  constructor() {
  }

  ngOnInit() {
    this.term = new Terminal({
      fontFamily: '"Monaco", "Consolas", "monospace"',
      fontSize: 12,
      rightClickSelectsWord: true,
      theme: {
        background: '#1f1b1b'
      }
    });
  }

  ngAfterViewInit() {
    this.term.open(this.terminal.nativeElement);
    this.term.resize(100, 26);
    const wsUrl = `ws://${window.location.host}${this.url}`;
    this.ws = new WebSocket(wsUrl);
    this.ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      const message = data['message'];
      this.term.write(message);
    };
  }

  onConfirm(): void {
    this.closeModal.emit(true);
  }

  ngOnDestroy(): void {
    this.ws.close();
  }

}
