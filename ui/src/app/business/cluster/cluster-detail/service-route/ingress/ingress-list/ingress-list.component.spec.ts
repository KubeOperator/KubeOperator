import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { IngressListComponent } from './ingress-list.component';

describe('IngressListComponent', () => {
  let component: IngressListComponent;
  let fixture: ComponentFixture<IngressListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ IngressListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(IngressListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
