import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DaemonSetListComponent } from './daemon-set-list.component';

describe('DaemonSetListComponent', () => {
  let component: DaemonSetListComponent;
  let fixture: ComponentFixture<DaemonSetListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DaemonSetListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DaemonSetListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
