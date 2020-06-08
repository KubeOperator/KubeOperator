import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DaemonSetComponent } from './daemon-set.component';

describe('DaemonSetComponent', () => {
  let component: DaemonSetComponent;
  let fixture: ComponentFixture<DaemonSetComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DaemonSetComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DaemonSetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
