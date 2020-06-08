import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StatefulSetListComponent } from './stateful-set-list.component';

describe('StatefulSetListComponent', () => {
  let component: StatefulSetListComponent;
  let fixture: ComponentFixture<StatefulSetListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StatefulSetListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatefulSetListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
